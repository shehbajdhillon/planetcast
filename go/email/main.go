package email

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"go.uber.org/zap"
)

type EmailConnectProps struct {
	Logger *zap.Logger
}

type Email struct {
	ses    *ses.SES
	logger *zap.Logger
}

func Connect(args EmailConnectProps) *Email {
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION := os.Getenv("AWS_REGION")

	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
	})
	sesClient := ses.New(session)

	if err != nil {
		panic(err)
	}

	args.Logger.Info("Email Client Started")

	return &Email{ses: sesClient, logger: args.Logger}

}

type sendEmailProps struct {
	Sender    string
	Recipient string
	Subject   string
	HtmlBody  string
	TextBody  string
}

func (e *Email) sendEmail(args sendEmailProps) error {

	charSet := "UTF-8"

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses:  []*string{aws.String(args.Recipient)},
			CcAddresses:  []*string{},
			BccAddresses: []*string{aws.String("shehbajdhillon0@gmail.com")},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(args.HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(args.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(args.Subject),
			},
		},
		Source: aws.String(args.Sender),
	}

	result, err := e.ses.SendEmail(input)

	if err == nil {
		e.logger.Info("Email sent to address", zap.String("recipient", args.Recipient), zap.String("email_output", result.GoString()))
	} else {
		e.logger.Error("Could not send email to address", zap.String("recipient", args.Recipient), zap.Error(err))
	}

	return err

}

func getHTMLTemplate(args interface{}, filePath string) (string, error) {
	var templateBuffer bytes.Buffer

	dir, err := filepath.Abs(filePath)

	if err != nil {
		return "", err
	}

	path := filepath.Join(dir)
	htmlData, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))
	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", args)

	if err != nil {
		return "", err
	}

	return templateBuffer.String(), nil

}

type DubbingAlertProps struct {
	TargetLanguage string
	TeamSlug       string
	ProjectId      int64
	UserEmail      string
}

func (e *Email) DubbingStartAlert(args DubbingAlertProps) {
	html := fmt.Sprintf(
		"Hey!<br /><br />We have started processing %s Dubbing for your project. We will alert you when it's done.<br /><br />You can access your project through this link:<br /><br />www.planetcast.ai/%s/%d<br /><br />Thanks,<br />Team PlanetCast",
		args.TargetLanguage, args.TeamSlug, args.ProjectId,
	)
	sendEmailArgs := sendEmailProps{
		Sender:    "PlanetCast <hello@planetcast.ai>",
		Recipient: args.UserEmail,
		HtmlBody:  html,
		TextBody:  html,
		Subject:   fmt.Sprintf("%s Dubbing Started", args.TargetLanguage),
	}
	e.sendEmail(sendEmailArgs)
}

func (e *Email) DubbingEndedAlert(args DubbingAlertProps) {
	html := fmt.Sprintf(
		"Hey!<br /><br />We have finished processing %s Dubbing for your project.<br /><br />You can access it through your project link:<br /><br />www.planetcast.ai/%s/%d<br /><br />Thanks,<br />Team PlanetCast",
		args.TargetLanguage, args.TeamSlug, args.ProjectId,
	)
	sendEmailArgs := sendEmailProps{
		Sender:    "PlanetCast <hello@planetcast.ai>",
		Recipient: args.UserEmail,
		HtmlBody:  html,
		TextBody:  html,
		Subject:   fmt.Sprintf("%s Dubbing Processed", args.TargetLanguage),
	}
	e.sendEmail(sendEmailArgs)
}
