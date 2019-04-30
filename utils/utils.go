package utils

import (
	"archive/zip"
	"bytes"
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	regular = `^1([38][0-9]|14[57]|5[^4])\d{8}$`
)

func UUID() string {
	u1, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return u1.String()
}

func GenerateSID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func Max(num1, num2 int) int {
	if num1 > num2 {
		return num1
	}
	return num2
}

func MaxInt32(num1, num2 int32) int32 {
	if num1 > num2 {
		return num1
	}
	return num2
}

func MaxInt64(num1, num2 int64) int64 {
	if num1 > num2 {
		return num1
	}
	return num2
}

func Min(num1, num2 int) int {
	if num1 < num2 {
		return num1
	}
	return num2
}

func MinInt32(num1, num2 int32) int32 {
	if num1 < num2 {
		return num1
	}
	return num2
}

func MinInt64(num1, num2 int64) int64 {
	if num1 < num2 {
		return num1
	}
	return num2
}

func Maxmin(num1, num2 int) (max, min int) {
	if num1 < num2 {
		return num2, num1
	}
	return num1, num2
}

func MaxminInt32(num1, num2 int32) (max, min int32) {
	if num1 < num2 {
		return num2, num1
	}
	return num1, num2
}

func MaxminInt64(num1, num2 int64) (max, min int64) {
	if num1 < num2 {
		return num2, num1
	}
	return num1, num2
}

func Swap(x, y int64) (int64, int64) {
	return y, x
}

func SwapString(x, y string) (string, string) {
	return y, x
}

func RandomString(c int) string {
	bs := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	result := []byte{}
	// result := make([]byte, c)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < c; i++ {
		result = append(result, bs[r.Intn(len(bs))])
		// result[i] = bytes[r.Intn(len(bytes))]
	}

	return string(result)
}

func RandomIntger(c int) int32 {
	if c <= 0 {
		return -1
	}

	bs := []byte("0123456789")

	result := make([]byte, c)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < c; i++ {
		result[i] = bs[r.Intn(len(bs))]
	}

	v, err := strconv.Atoi(string(result))
	if err != nil {
		return -1
	}

	return int32(v)
}

func CompareMapStringString(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		if v1, ok := m2[k]; !ok || v1 != v {
			return false
		}
	}

	return true
}

func ValidatePhone(mobileNum string) bool {
	if len(mobileNum) == 0 {
		return false
	}

	reg := regexp.MustCompile(regular)

	ok := reg.MatchString(mobileNum)

	return ok
}

func ValidateEmail(email string) bool {
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		return false
	}
	return true
}

func ValidateIDCard(card string) bool {
	//验证15位身份证，15位的是全部数字
	if m, _ := regexp.MatchString(`^(\d{15})$`, card); !m {
		return false
	}

	//验证18位身份证，18位前17位为数字，最后一位是校验位，可能为数字或字符X。
	if m, _ := regexp.MatchString(`^(\d{17})([0-9]|X)$`, card); !m {
		return false
	}

	return true
}

func Substr(str string, start int, end int) (string, error) {
	rs := []byte(str)
	length := len(str)

	if start < 0 || start > length {
		return "", errors.New("start is wrong")
	}

	if end < 0 || end > length {
		return "", errors.New("end is wrong")
	}

	return string(rs[start:end]), nil
}

func HashCode(s string) uint32 {
	v := crc32.ChecksumIEEE([]byte(s))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
func HashCodes(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", HashCode(buf.String()))
}

func ZipDir(dir, zipFile string, f func(file string)) error {
	fz, err := os.Create(zipFile)
	if err != nil {
		return err
	}

	w := zip.NewWriter(fz)
	defer w.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			newPath := strings.Replace(path, "\\", "/", -1)
			var fDest io.Writer
			fDest, err = w.Create(newPath[len(dir)+1:])
			if err != nil {
				return err
			}

			if f != nil {
				f(newPath)
			}
			var fSrc *os.File
			fSrc, err = os.Open(newPath)
			if err != nil {
				return err
			}
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func UnzipDir(zipFile, dir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {
		err = func() error {
			path := dir + string(filepath.Separator) + f.Name
			_ = os.MkdirAll(filepath.Dir(path), 0755)
			var fDest *os.File
			fDest, err = os.Create(path)
			if err != nil {
				return err
			}
			defer fDest.Close()

			var fSrc io.ReadCloser
			fSrc, err = f.Open()
			if err != nil {
				return err
			}
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				return err
			}
			return nil
		}()

		if err != nil {

		}
	}

	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
