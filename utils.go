package main

func checkErr(message string, err error) {
	if err != nil {
		infoLabel.SetStyleName("error")
		if message != "" {
			infoLabel.SetText(message + ": " + err.Error())
		} else {
			infoLabel.SetText(err.Error())
		}
	}
}
