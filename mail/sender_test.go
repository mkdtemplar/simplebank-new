package mail

import "testing"

func TestGmailSender_SendEmail(t *testing.T) {
	type fields struct {
		name              string
		fromEmailAddress  string
		fromEmailPassword string
	}
	type args struct {
		subject     string
		content     string
		to          []string
		cc          []string
		bcc         []string
		attachFiles []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := &GmailSender{
				name:              tt.fields.name,
				fromEmailAddress:  tt.fields.fromEmailAddress,
				fromEmailPassword: tt.fields.fromEmailPassword,
			}
			if err := sender.SendEmail(tt.args.subject, tt.args.content, tt.args.to, tt.args.cc, tt.args.bcc, tt.args.attachFiles); (err != nil) != tt.wantErr {
				t.Errorf("GmailSender.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
