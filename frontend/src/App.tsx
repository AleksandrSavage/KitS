
import { useEffect, useState } from "react";
import "./styles/App.css";

import {
    SelectAudioFile,
    Play,
    Pause,
    Resume,
    Stop,
    GetState,
} from "../wailsjs/go/main/App.js";

import {
    EventsOn,
    EventsOff,
} from "../wailsjs/runtime/runtime.js";

type PlayerEvent = {
    type?: string;
    state?: string;
    error?: string;
};

const App = ()=> {
    const [filePath, setFilePath] = useState("");
    const [state, setState] = useState("idle");
    const [error, setError] = useState("");
    const [busy, setBusy] = useState(false);

    // Подписываемся на события от Go.
    useEffect(() => {
        EventsOn("state_changed", (event: PlayerEvent) => {
            setState(event.state ?? "unknown");
        });

        EventsOn("error", (event: PlayerEvent) => {
            setError(event.error ?? "Unknown error");
        });

        EventsOn("track_ended", () => {
            setState("ended");
        });

        // Получаем начальное состояние плеера.
        GetState()
            .then(setState)
            .catch((err) => setError(String(err)));

        // Убираем обработчики при размонтировании компонента.
        return () => {
            EventsOff(
                "state_changed",
                "error",
                "track_ended",
            );
        };
    }, []);

    // Универсальная обработка команд плеера.
    async function runCommand(
        command: () => Promise<void>,
    ) {
        setBusy(true);
        setError("");

        try {
            await command();
        } catch (err) {
            setError(String(err));
        } finally {
            setBusy(false);
        }
    }

    // Открытие системного окна выбора файла.
    async function chooseFile() {
        setError("");

        try {
            const path = await SelectAudioFile();

            // Пустая строка означает, что пользователь отменил выбор.
            if (!path) return;

            setFilePath(path);
        } catch (err) {
            setError(String(err));
        }
    }

    // Запуск выбранного файла.
    function play() {
        if (!filePath) {
            setError("Сначала выбери аудиофайл");
            return;
        }

        return runCommand(() => Play(filePath));
    }

    const fileName =
        filePath.split(/[\\/]/).pop() || "Файл не выбран";

    return (
        <main className="player">
            <h1>Audio Player</h1>

            <section className="track">
                <p className="label">Текущий файл</p>
                <p className="filename" title={filePath}>
                    {fileName}
                </p>

                <button
                    onClick={chooseFile}
                    disabled={busy}
                >
                    Выбрать файл
                </button>
            </section>

            <section className="controls">
                <button
                    onClick={play}
                    disabled={busy || !filePath}
                >
                    Play
                </button>

                <button
                    onClick={() => runCommand(() => Pause())}
                    disabled={busy || state !== "playing"}
                >
                    Pause
                </button>

                <button
                    onClick={() => runCommand(() => Resume())}
                    disabled={busy || state !== "paused"}
                >
                    Resume
                </button>

                <button
                    onClick={() => runCommand(() => Stop())}
                    disabled={
                        busy ||
                        (state !== "playing" && state !== "paused")
                    }
                >
                    Stop
                </button>
            </section>

            <section className="status">
                <p className="label">Состояние плеера</p>

                <p className="state">
                    {state.toUpperCase()}
                </p>
            </section>

            {busy && (
                <p className="info">
                    Выполняется команда...
                </p>
            )}

            {error && (
                <section className="error">
                    <p>Ошибка</p>
                    <span>{error}</span>
                    <button onClick={() => setError("")}>
                        Закрыть
                    </button>
                </section>
            )}
        </main>
    );
}

export default App;
