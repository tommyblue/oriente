import React from "react";
import "./game.css";
import back from "./assets/back.png";

function Game({ gameState }) {
    return (
        <div className="game-space">
            {" "}
            <img src={back} alt="Deck" className="cover deck" />
            <GamePrize num={gameState.prize_cards} />
        </div>
    );
}

function GamePrize({ num }) {
    let shiftPx = 4;
    const deck = [];
    for (let i = 0; i < num; i++) {
        const shift = shiftPx * i;
        console.log("prize");
        deck.push(
            <img
                key={`prize_${i}`}
                src={back}
                alt="Deck"
                className={`cover prize ${i === 0 ? "" : "over"}`}
                style={{ top: shift, left: shift }}
            />
        );
    }
    return <div className="prize-deck">{deck}</div>;
}

export default Game;
