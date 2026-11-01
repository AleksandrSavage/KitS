import { PlayAudio, StopAudio, PauseAudio, ResumeAudio } from "../wailsjs/go/main/App.js";

const App = ()=> {
const play = async () => {
    await ResumeAudio();
}

const init = async () => {
    await PlayAudio("test/audio-test.mp3");
};

const stop = async () => {
    await PauseAudio();
};    


return (
    <div>
        <button onClick={init}>bitton</button>
        <button onClick={play}>play</button>
        <button onClick={stop}>stop bitton</button>
    </div>
)
}

export default App
