import { PlayAudio, StopAudio } from "../wailsjs/go/main/App.js";

const App = ()=> {

const play = async () => {
        PlayAudio("test/audio-test.mp3");
    };

const stop = async () => {
       StopAudio();
    };    

return (
    <div>
        <button onClick={play}>bitton</button>
        <button onClick={stop}>stop bitton</button>
    </div>
)
}

export default App
