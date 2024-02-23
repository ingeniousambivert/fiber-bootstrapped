package modules

import (
	"fmt"
)

type Mailer struct {
	From string
}

func (m *Mailer) Send(params map[string]interface{}) {
	to, _ := params["to"].(string)
	subject, _ := params["subject"].(string)
	body, _ := params["body"].(string)

	fmt.Print("Mailer Module - Start\n")
	fmt.Printf("To - %s\n", to)
	fmt.Printf("Subject - %s\n", subject)
	fmt.Printf("Body - %s\n", body)
	fmt.Print("Mailer Module - End\n")
}
