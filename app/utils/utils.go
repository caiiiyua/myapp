package utils

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/nu7hatch/gouuid"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

// assert
func AssertNoError(err error, msg string) {
	if err != nil {
		Assert(false, msg+"("+err.Error()+")")
	}
}

func Assert(assertion bool, msg string) {
	if !assertion {
		err := "assertion failed"
		if msg != "" {
			err += ": " + msg
		}
		panic(err)
	}
}

// time.layout
const (
	DefaultDateTimeFull = "2016-01-02 15:04:05.999"
	DefaultDateTime     = "2016-01-02 15:04:05"
	DefaultDate         = "2016-01-02"
	DefaultTime         = "15:04:05"
)

// time.format
// FormatDefault returns time as DefaultDateTime format
func FormatDefault(t time.Time) string {
	return time.Time.Format(t, DefaultDateTime)
}

type RichTime time.Time

// Yesterday for short
func (r RichTime) Yesterday() time.Time {
	var t = time.Time(r)
	return t.AddDate(0, 0, -1)
}

// revel
// Get config from revel, panic if not exists
func ForceGetConfig(key string) string {
	v, exists := revel.Config.String(key)
	Assert(exists, fmt.Sprintf("Missing revel app config for key:%s", key))
	// fmt.Printf("[%s]: %v", key, v)
	return v
}

// uuid
func Uuid() string {
	u4, err := uuid.NewV4()
	AssertNoError(err, "")

	return u4.String()
}

// SHA1
func Sha1(content string) string {
	h := sha1.New()
	io.WriteString(h, content)

	return fmt.Sprintf("%x", h.Sum(nil))
}

// 随机字符串
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randInt(48, 57))
	}

	return string(bytes)
}

func randInt(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return min + r.Intn(max-min)
}

// 生成并保存缩略图
func MakeAndSaveThumbnail(fromFile string, toFile string, w, h int) error {
	tnImage, err := MakeThumbnail(fromFile, w, h)
	if err != nil {
		return err
	}
	return imaging.Save(tnImage, toFile)
}

// 生成并保存缩略图
func MakeAndSaveThumbnailFromReader(reader io.Reader, toFile string, w, h int) error {
	tnImage, err := MakeThumbnailFromReader(reader, w, h)
	if err != nil {
		return err
	}
	return imaging.Save(tnImage, toFile)
}

func MakeAndSaveFromReader(reader io.Reader, toFile string, t string, w, h int) error {
	tnImage, err := MakeFromReader(reader, t, w, h)
	if err != nil {
		return err
	}
	return imaging.Save(tnImage, toFile)
}

// 生成缩略图
func MakeThumbnail(fromFile string, w, h int) (image *image.NRGBA, err error) {
	srcImage, err := imaging.Open(fromFile)
	if err != nil {
		return nil, err
	}

	image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)
	return
}

func MakeThumbnailFromReader(reader io.Reader, w, h int) (image *image.NRGBA, err error) {
	srcImage, err := Open(reader)
	if err != nil {
		return nil, err
	}

	image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)
	return
}

func MakeFromReader(reader io.Reader, t string, w, h int) (image *image.NRGBA, err error) {
	srcImage, err := Open(reader)
	if err != nil {
		return nil, err
	}

	switch t {
	case "thumbnail":
		image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)
	case "resize":
		image = imaging.Resize(srcImage, w, h, imaging.Lanczos)
	case "fit":
		image = imaging.Fit(srcImage, w, h, imaging.Lanczos)
	default:
		image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)
	}

	return
}

func MakeFromImage(srcImage image.Image, t string, w, h int) (image *image.NRGBA, err error) {
	switch t {
	case "thumbnail":
		image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)
	case "resize":
		image = imaging.Resize(srcImage, w, h, imaging.Lanczos)
	case "fit":
		image = imaging.Fit(srcImage, w, h, imaging.Lanczos)
	default:
		image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)
	}

	return image, nil
}

func MakeAndSaveFromReaderMax(reader io.Reader, toFile string, w, h int) error {
	tnImage, err := MakeFromReaderMax(reader, w, h)
	if err != nil {
		return err
	}
	return imaging.Save(tnImage, toFile)
}

func MakeFromReaderMax(reader io.Reader, maxW, maxH int) (image *image.NRGBA, err error) {
	srcImage, err := Open(reader)
	if err != nil {
		return nil, err
	}

	srcBounds := srcImage.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	if w > maxW {
		w = maxW
	}
	if h > maxH {
		h = maxH
	}

	image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)

	return
}

func MakeAndSaveFromReaderMaxWithMode(reader io.Reader, t string, toFile string, w, h int) error {
	tnImage, err := MakeFromReaderMaxWithMode(reader, t, w, h)
	if err != nil {
		return err
	}
	return imaging.Save(tnImage, toFile)
}

func MakeFromReaderMaxWithMode(reader io.Reader, t string, maxW, maxH int) (image *image.NRGBA, err error) {
	srcImage, err := Open(reader)
	if err != nil {
		return nil, err
	}

	srcBounds := srcImage.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	if w > maxW {
		w = maxW
	}
	if h > maxH {
		h = maxH
	}

	image, _ = MakeFromImage(srcImage, t, w, h)

	return
}

// Open loads an image from file
func Open(reader io.Reader) (img image.Image, err error) {
	img, _, err = image.Decode(reader)
	if err != nil {
		return
	}

	img = toNRGBA(img)
	return
}

// This function used internally to convert any image type to NRGBA if needed.
func toNRGBA(img image.Image) *image.NRGBA {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*image.NRGBA); ok {
			return src0
		}
	}
	return imaging.Clone(img)
}

func ToJSON(o interface{}) string {
	b, err := json.Marshal(o)
	AssertNoError(err, "ToJSON")

	return string(b)
}

func FromJSON(s string, o interface{}) {
	err := json.Unmarshal([]byte(s), o)
	AssertNoError(err, "FromJSON)")
}

//panicable
type CacheDataLoader func(string) interface{}

func Cache(key string, target interface{}, loader CacheDataLoader) {
	CacheWithExpires(key, target, loader, cache.FOREVER)
}

var cacheKeys []string

func GetCacheKeys() []string {
	return cacheKeys
}

func CacheWithExpires(key string, target interface{}, loader CacheDataLoader, expires time.Duration) {
	if err := cache.Get(key, target); err != nil {
		values := loader(key)
		setValueToAddress(target, values)
		cacheKeys = append(cacheKeys, key)
		go cache.Set(key, values, expires)
	}
}

func ClearCache(key string) {
	go cache.Delete(key)
}

func setValueToAddress(target interface{}, value interface{}) {
	p := reflect.ValueOf(target)
	Assert(p.Type().Kind() == reflect.Ptr, "target should be Pointer")

	v := p.Elem()

	Assert(v.CanSet(), "target should be CanSet")
	v.Set(reflect.ValueOf(value))
}

//
func TypeOfTarget(v interface{}) (typ reflect.Type) {
	typ = reflect.TypeOf(v)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return
}

func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
