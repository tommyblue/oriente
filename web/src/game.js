import React from "react";
import "./game.css";
import back from "./assets/back.jpg";

function Game(props) {
    return (
        <div>
            {" "}
            <img src={back} alt="Deck" className="cover" />
        </div>
    );
}

export default Game;
