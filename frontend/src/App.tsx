
import { useEffect, useState } from "react";
import "./styles/App.css";

import {
    SelectAudioFile,
    Play,
    Pause,
    Resume,
    Stop,
    GetState,
    GetProgress,
    Seek,
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

    const [isSeeking, setIsSeeking] = useState(false);
    const [seekValue, setSeekValue] = useState(0);

    function formatTime(ms: number): string {
        const totalSeconds = Math.max(0, Math.floor(ms / 1000));
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;

        return `${minutes}:${String(seconds).padStart(2, "0")}`;
    }

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

    async function chooseFile() {
        setError("");

        try {
            const path = await SelectAudioFile();

            if (!path) return;

            if (
                path !== activePath &&
                (state === "playing" || state === "paused")
            ) {
                await Stop();
            }

            setFilePath(path);

            if (path !== activePath) {
                setActivePath("");
                setProgress({
                    positionMs: 0,
                    durationMs: 0,
                });
            }
        } catch (err) {
            setError(String(err));
        }
    }

    async function togglePlayback() {
        if (!filePath) {
            setError("Сначала выбери аудиофайл");
            return;
        }

        await runCommand(async () => {
            if (filePath !== activePath) {
                await Play(filePath);
                setActivePath(filePath);
                return;
            }

            if (state === "playing") {
                await Pause();
                return;
            }

            if (state === "paused") {
                await Resume();
                return;
            }

            await Play(filePath);
            setActivePath(filePath);
        });
    }

    async function commitSeek(positionMs: number) {
        const target = Math.max(
            0,
            Math.min(positionMs, visibleProgress.durationMs)
        );

        await runCommand(async () => {
            await Seek(target);

            setProgress((prev) => ({
                ...prev,
                positionMs: target,
            }));
        });

        setIsSeeking(false);
    }

    const fileName =
        filePath.split(/[\\/]/).pop() || "Файл не выбран";

    const visibleProgress =
        activePath === filePath
            ? progress
            : { positionMs: 0, durationMs: 0 };

    const displayedPositionMs = isSeeking
        ? seekValue
        : visibleProgress.positionMs;

    const progressPercent =
        visibleProgress.durationMs > 0
            ? Math.min(
                100,
                Math.max(
                    0,
                    (displayedPositionMs / visibleProgress.durationMs) * 100
                )
            )
            : 0;

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
                    className="play-button"
                    onClick={togglePlayback}
                    disabled={busy || !filePath}
                >
                    {state === "playing" && activePath === filePath
                        ? "Pause"
                        : "Play"}
                </button>
            </section>

            <section className="timeline">
                <span className="time-current">
                    {formatTime(displayedPositionMs)}
                </span>

                <input
                    className="progress-slider"
                    type="range"
                    min={0}
                    max={visibleProgress.durationMs}
                    step={100}
                    value={Math.min(
                        displayedPositionMs,
                        visibleProgress.durationMs
                    )}
                    disabled={
                        busy ||
                        visibleProgress.durationMs <= 0 ||
                        (state !== "playing" && state !== "paused")
                    }
                    style={{
                        background: `linear-gradient(
                            to right,
                            #66cf9c ${progressPercent}%,
                            #414653 ${progressPercent}%
                        )`,
                    }}
                    onPointerDown={() => {
                        setIsSeeking(true);
                        setSeekValue(visibleProgress.positionMs);
                    }}
                    onChange={(event) => {
                        setIsSeeking(true);
                        setSeekValue(Number(event.target.value));
                    }}
                    onPointerUp={(event) => {
                        void commitSeek(
                            Number(event.currentTarget.value)
                        );
                    }}
                    onKeyUp={(event) => {
                        if (
                            [
                                "ArrowLeft",
                                "ArrowRight",
                                "Home",
                                "End",
                                "PageUp",
                                "PageDown",
                            ].includes(event.key)
                        ) {
                            void commitSeek(
                                Number(event.currentTarget.value)
                            );
                        }
                    }}
                />

                <span className="time-duration">
                    {formatTime(visibleProgress.durationMs)}
                </span>
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
