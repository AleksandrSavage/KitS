import { PlayAudio } from "../wailsjs/go/main/App.js";

const App = ()=> {

const handleClick = async () => {
        PlayAudio("test/audio-test.mp3");
    };

return (
    <div>
        <button onClick={handleClick}>bitton</button>
    </div>
)
}

export default App
