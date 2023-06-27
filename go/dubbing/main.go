package dubbing

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"planetcastdev/database"
)

type WhisperOutput struct {
	Text     string    `json:"text"`
	Language string    `json:"language"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	Id    int64   `json:"id"`
	Seek  int64   `json:"seek"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

func Dub(sourceLanguage database.SupportedLanguage, targetLanguage database.SupportedLanguage, fileNameIdentifier string, file io.ReadSeeker) {

	file.Seek(0, io.SeekStart)

	//write file to disk
	body, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
	}
	err = ioutil.WriteFile("./"+fileNameIdentifier+".mp4", body, 0644)
	if err != nil {
		log.Println(err.Error())
	}

	//run whisper
	cmd := exec.Command("whisper", "./"+fileNameIdentifier+".mp4", "--model", "tiny", "--output_format", "json")
	log.Println("Running Whisper:", cmd.Args)

	body, err = cmd.Output()
	if err != nil {
		log.Println("Whisper error:", err.Error())
	} else {
		log.Println("Whisper output file:", fileNameIdentifier+".json")
	}

	var whisperOutput WhisperOutput
	whisperJson, _ := ioutil.ReadFile("./" + fileNameIdentifier + ".json")
	json.Unmarshal(whisperJson, &whisperOutput)

	cmd = exec.Command("rm", "./"+fileNameIdentifier+".mp4", "./"+fileNameIdentifier+".json")
	log.Println(cmd.Args)
	body, err = cmd.Output()

	log.Println(whisperOutput.Segments)

	//get whisper output and parse it
}
