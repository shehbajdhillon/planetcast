package dubbing

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"planetcastdev/database"

	"github.com/tabbed/pqtype"
)

type WhisperOutput struct {
	Language string    `json:"language"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	Id    int64   `json:"id"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

func getTranscript(fileNameIdentifier string, file io.ReadSeeker) WhisperOutput {

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
	cmd := exec.Command("whisper", "./"+fileNameIdentifier+".mp4", "--model", "medium", "--output_format", "json")
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

	return whisperOutput
}

func CreateTransformation(ctx context.Context, projectId int64, targetLanguage database.SupportedLanguage, fileNameIdentifier string, file io.ReadSeeker, queries *database.Queries, isSource bool) (database.Transformation, error) {

	transcript := getTranscript(fileNameIdentifier, file)
	jsonBytes, err := json.Marshal(transcript)

	transformation, err := queries.CreateTransformation(ctx, database.CreateTransformationParams{
		ProjectID:      projectId,
		TargetLanguage: targetLanguage,
		TargetMedia:    fileNameIdentifier + ".mp4",
		Transcript:     pqtype.NullRawMessage{RawMessage: jsonBytes, Valid: true},
		IsSource:       isSource,
	})

	if err != nil {
		log.Println("Error occured:", err.Error())
		return database.Transformation{}, err
	}

	return transformation, nil
}
