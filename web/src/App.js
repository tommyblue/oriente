import React, { useState, useEffect } from "react";
import Intro from "./intro";
import Game from "./game";
import "./App.css";

function App() {
    const [game, setGame] = useState(null);
    const [player, setPlayer] = useState(null);
    const [started, setStarted] = useState(false);
    const [timer, setTimer] = useState(null);
    const [gameState, setGameState] = useState({});

    useEffect(() => {
        if (!player || !game) {
            return;
        }
        if (timer !== null) {
            return;
        }
        const t = setInterval(() => {
            if (!game || !player) {
                return;
            }
            fetch(`http://localhost:8000/game/${game}/${player}`)
                .then(res => {
                    if (!res.ok) {
                        throw res.statusText; // TODO: manage the error
                    }
                    return res.json();
                })
                .then(res => {
                    setGameState(res);
                    if (res.game_started === true) {
                        setStarted(true);
                        clearInterval(t);
                    }
                })
                .catch(err => console.error(err));
        }, 3000);
        setTimer(t);
        return () => {
            if (timer !== null) {
                console.log("cleanup");
            }
        };
    }, [game, player, started, timer]);

    return (
        <div className="App">
            <Debug game={game} player={player} started={started} />
            {game && player && started ? (
                <Game />
            ) : (
                <Intro
                    setGame={setGame}
                    game={game}
                    setPlayer={setPlayer}
                    player={player}
                    waitingPlayers={timer !== null}
                    gameState={gameState}
                />
            )}
        </div>
    );
}

function Debug({ game, player, started }) {
    return (
        <div className="debug">
            <div>Current game: {game}</div>
            <div>Current player: {player}</div>
            <div>Game started: {started ? "true" : "false"}</div>
            {game == null || player === null ? (
                <span />
            ) : (
                <a
                    href={`http://localhost:8000/game/${game}/${player}`}
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Debug
                </a>
            )}
        </div>
    );
}

export default App;
