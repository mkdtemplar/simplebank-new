package mail

import (
	"testing"

	"github.com/mkdtemplar/simplebank-new/util"
	"github.com/stretchr/testify/require"
)

func TestGmailSender_SendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test email"
	content := `<h1>This is the test email for simplebank</h1>`
	to := []string{"sagitariusim@live.com"}
	attachFiles := []string{"../README.MD"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
