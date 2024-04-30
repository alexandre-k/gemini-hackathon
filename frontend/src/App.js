import React from "react";
import axios from 'axios';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import { Brain, Translate } from "@phosphor-icons/react";
import IconButton from '@mui/material/IconButton';
import Speech from 'speak-tts' // es6
// import { useQuery } from '@tanstack/react-query'

import './App.css';
import Webcam from "react-webcam";


function App() {
    const [isLoading, setIsLoading] = React.useState(false)
    const dataUrlToBlob = dataUrl => {
        var byteString = atob(dataUrl.split(',')[1]);
        var mimeString = dataUrl.split(',')[0].split(':')[1].split(';')[0];
        var ab = new ArrayBuffer(byteString.length);
        var ia = new Uint8Array(ab);
        for (var i = 0; i < byteString.length; i++) {
            ia[i] = byteString.charCodeAt(i);
        }
        var blob = new Blob([ab], {type: mimeString});
        return blob;

    };
    const speech = new Speech();
    speech.init({"lang": "en-US"}).then((data) => {
        // The "data" object contains the list of available voices and the voice synthesis params
        console.log("Speech is ready, voices are available", data);
    }).catch(e => {
        console.error("An error occured while initializing : ", e);
    });
    const readOutLoud = (text) => {
        speech.speak({
            text,
        }).then(() => {
            console.log("Success !")
        }).catch(e => {
            console.error("An error occurred :", e)
        });
    };


    const generateAiResponse = async (action, imageBlob) => {
        readOutLoud(`"I will ${action} this, wait a sec!`);
        let formData = new FormData();
        formData.append("image", imageBlob);
        formData.append("imageType", imageBlob.type.split("/")[1]);

        const response = await axios(`${process.env.REACT_APP_API_URL}/gemini/${action}`, {
            method: 'post',
            data: formData
        });
        console.log(response.data.content);
        readOutLoud(response.data.content);
    };

    const webcamRef = React.useRef(null);
    const capture = React.useCallback(
        (action) => {
            setIsLoading(true);
            const imageDataUrl = webcamRef.current.getScreenshot();
            generateAiResponse(action, dataUrlToBlob(imageDataUrl));
            setIsLoading(false);
        },
        [webcamRef]
    );

    const videoConstraints = {
        facingMode: { exact: "environment" }
    };

    return (
        <Grid sx={{ backgroundColor: '#000'}} direction="column" justifyContent="center" spacing={0} alignItems="stretch">
          <Grid item xs={12} sx={{ height: '80vh', textAlign: 'center'}}>
            <Webcam
                sx={{ width: '100%', textAlign: 'center'}} 
                id="webcam"
                ref={webcamRef}
                audio={false}
                screenshotFormat="image/jpeg"
                videoConstraints={videoConstraints}
                onClick={()=> readOutLoud('Click on the bottom left to get an overview of what you camera sees, and on the bottom right to have an overview of what is written.')}
                />
          </Grid>
          <Grid container sx={{ height: '20vh'}}>
            <Grid item xs={6} sx={{ backgroundColor: 'lightGrey', textAlign: 'center', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }} item xs={6}>
              <IconButton
                color="primary"
                size="large"
                aria-label="imageDescription"
                color="success"
                disabled={isLoading}
                onClick={() => capture("describe")}>
                <Brain size={64} color="#fff" weight="fill" />
              </IconButton>
            </Grid>
            <Grid item xs={6} sx={{ backgroundColor: 'grey', textAlign: 'center', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center'  }} item xs={6}>
              <IconButton
                color="primary"
                size="large"
                aria-label="imageTranslation"
                color="success"
                disabled={isLoading}
                onClick={() => capture("translate")}>
                <Translate size={64} color="#fff" weight="fill" />
              </IconButton>
            </Grid>
          </Grid>
        </Grid>
  );
}

export default App;
