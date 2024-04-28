import React from "react";
import axios from 'axios';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import { Brain, Translate } from "@phosphor-icons/react";
import IconButton from '@mui/material/IconButton';
import Speech from 'speak-tts' // es6
// import { useQuery } from '@tanstack/react-query'

import './App.css';
import Webcam from "react-webcam";


function App() {
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
    speech.init().then((data) => {
        // The "data" object contains the list of available voices and the voice synthesis params
        console.log("Speech is ready, voices are available", data);
    }).catch(e => {
        console.error("An error occured while initializing : ", e);
    });
    const readOutloud = (text) => {
        speech.speak({
            text,
        }).then(() => {
            console.log("Success !")
        }).catch(e => {
            console.error("An error occurred :", e)
        });
    };


    const generateAiResponse = async (action, imageBlob) => {
        readOutloud("Wait a second...");
        let formData = new FormData();
        formData.append("image", imageBlob);
        formData.append("imageType", imageBlob.type.split("/")[1]);

        const response = await axios(`${process.env.REACT_APP_API_URL}/gemini/${action}`, {
            method: 'post',
            data: formData
        });
        console.log(response.data.content);
        readOutloud(response.data.content);
    };

    // useEffect(() => {
    //     generateAiDescription()
    // }, []);
    // const { isPending, isError, data, error } = useQuery({
    //     queryKey: [],
    //     queryFn: generateAiDescription,
    // })

    // if (isPending) {
    //     return <span>Loading...</span>
    // }

    // if (isError) {
    //     return <span>Error: {error.message}</span>
    // }

    const webcamRef = React.useRef(null);
    const capture = React.useCallback(
        (action) => {
            const imageDataUrl = webcamRef.current.getScreenshot();
            generateAiResponse(action, dataUrlToBlob(imageDataUrl));
        },
        [webcamRef]
    );

    const videoConstraints = {
        facingMode: { exact: "environment" }
    };

    return (
        <Container sx={{ bgcolor: '#000', display: 'flex', flexFlow: 'column', alignItems: 'center' }}>
        <Box sx={{ bgcolor: '#000', height: '100vh', width: '100vh', display: 'flex', alignItems: 'center' }}>
	        <Webcam
            id="webcam"
            ref={webcamRef}
            audio={false}
            screenshotFormat="image/jpeg"
            videoConstraints={videoConstraints}
          />
        </Box>
        <Box sx={{ display: 'flex', justifyContent: 'flex-center', position: 'absolute', bottom: '100px' }}>
          <IconButton color="primary" size="large" aria-label="fingerprint" color="success" onClick={() => capture("describe")}>
              <Brain size={64} color="#fff" weight="fill" />
            </IconButton>
          <IconButton color="primary" size="large" aria-label="fingerprint" color="success" onClick={() => capture("translate")}>
              <Translate size={64} color="#fff" weight="fill" />
            </IconButton>
        </Box>
      </Container>
  );
}

export default App;
