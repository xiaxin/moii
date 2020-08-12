package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)


/**
	TODO
		- 参数优化
		- 功能优化
		- 自定义变量
 */

/**
	生成字符串验证码
	返回数据
		answer  验证码字符串
		b64     图片 base64 编码
		err     错误
 */
func GenerateStringCaptcha(height, width, length int) (answer string, b64s string, err error) {
	driverString := base64Captcha.NewDriverString(height, width, 5, 0, length, "1234567890qwertyuioplkjhgfdsazxcvbnm", &color.RGBA{R: 0, G: 1, B: 0, A: 0}, []string{"wqy-microhei.ttc"})

	driver := driverString.ConvertFonts()

	c := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)

	_, content, answer := c.Driver.GenerateIdQuestionAnswer()

	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}

	b64s = item.EncodeB64string()

	return answer, b64s, err
}

/**
	生成字符串验证码
	返回数据
		answer  验证码字符串
		b64     图片 base64 编码
		err     错误
*/
func GenerateMathCaptcha(height, width int) (answer string, b64s string, err error)  {
	driverMath := base64Captcha.NewDriverMath(height, width, 5, 0, &color.RGBA{R: 0, G: 1, B: 0, A: 0}, []string{"wqy-microhei.ttc"} )

	driver := driverMath.ConvertFonts()

	c := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)

	_, content, answer := c.Driver.GenerateIdQuestionAnswer()

	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}

	b64s = item.EncodeB64string()

	return answer, b64s, err
}
