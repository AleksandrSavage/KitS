
import { useEffect, useState } from "react";
import "./styles/App.css";

import {
    SelectAudioFile,
    Play,
    Pause,
    Resume,
    GetState,
    GetProgress,
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

function App() {
    const [activePath, setActivePath] = useState("");
    const [filePath, setFilePath] = useState("");
    const [state, setState] = useState("idle");
    const [error, setError] = useState("");
    const [busy, setBusy] = useState(false);

    const [progress, setProgress] = useState({
        positionMs: 0,
        durationMs: 0,
    });

    // Форматирование времени: 65000 -> 1:05
    function formatTime(ms: number): string {
        const totalSeconds = Math.max(0, Math.floor(ms / 1000));

        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;

        return `${minutes}:${String(seconds).padStart(2, "0")}`;
    }

    // Получение прогресса воспроизведения.
    useEffect(() => {
        let cancelled = false;

        async function updateProgress() {
            try {
                const result = await GetProgress();

                if (!cancelled) {
                    setProgress(result);
                }
            } catch (err) {
                console.error("Progress error:", err);
            }
        }

        updateProgress();

        const timer = window.setInterval(updateProgress, 250);

        return () => {
            cancelled = true;
            window.clearInterval(timer);
        };
    }, []);

    // Подписка на события от Go.
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

        // Получаем начальное состояние.
        GetState()
            .then(setState)
            .catch((err) => setError(String(err)));

        return () => {
            EventsOff(
                "state_changed",
                "error",
                "track_ended"
            );
        };
    }, []);

    // Универсальная обработка команд плеера.
    async function runCommand(
        command: () => Promise<void>
    ): Promise<void> {
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

            if (!path) return;

            setFilePath(path);
        } catch (err) {
            setError(String(err));
        }
    }

    // Единая кнопка Play / Pause.
    async function togglePlayback() {
        if (!filePath) {
            setError("Сначала выбери аудиофайл");
            return;
        }

        await runCommand(async () => {
            // Выбран другой файл.
            if (filePath !== activePath) {
                await Play(filePath);
                setActivePath(filePath);
                return;
            }

            // Музыка играет — ставим на паузу.
            if (state === "playing") {
                await Pause();
                return;
            }

            // Музыка на паузе — продолжаем.
            if (state === "paused") {
                await Resume();
                return;
            }

            // Трек закончился или остановлен — запускаем заново.
            await Play(filePath);
            setActivePath(filePath);
        });
    }

    const fileName =
        filePath.split(/[\\/]/).pop() || "Файл не выбран";
        
    const visibleProgress =
        activePath === filePath
            ? progress
            : { positionMs: 0, durationMs: 0 };

    const progressPercent =
        visibleProgress.durationMs > 0
            ? Math.min(
                100,
                (visibleProgress.positionMs / visibleProgress.durationMs) * 100
            )
            : 0;

    return (
        <main className="player">
            <h1>Audio Player</h1>

            {/* Информация о треке */}
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

            {/* Управление воспроизведением */}
            <section className="controls">
                <button
                    className="play-button"
                    onClick={togglePlayback}
                    disabled={busy || !filePath}
                >
                    {state === "playing" && activePath === filePath
                        ? "Pause"
                        : "Play"}
                </button>
            </section>

            {/* Время и прогресс-бар */}
            <section className="timeline">
                <span className="time-current">
                    {formatTime(visibleProgress.positionMs)}
                </span>

                <div
                    className="progress-track"
                    role="progressbar"
                    aria-valuemin={0}
                    aria-valuemax={100}
                    aria-valuenow={progressPercent}
                >
                    <div
                        className="progress-fill"
                        style={{
                            width: `${progressPercent}%`,
                        }}
                    />
                </div>

                <span className="time-duration">
                    {formatTime(visibleProgress.durationMs)}
                </span>
            </section>

            {/* Состояние плеера */}
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

            {}
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
