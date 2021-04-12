function updateBoardState(board, selectionDiv, roomActions, playerName, playerID, roomID, state) {
  applyBoardState(board, state, roomID);
  applyCardStyle(selectionDiv, state, playerID);

  if (playerName) {
    roomActions.style.display = 'block';
  }
}

function applyBoardState(board, state, roomID) {
  const listPLayers = getCardBoardState(state);
  const roomStateString = getBoardStateString(state.RoomState);

  const boardInner = `<div>
    <span>board status : ${roomStateString}</span>
    ${listPLayers}
    <span>room id: ${roomID}</span>
  </div>`
  board.innerHTML = boardInner;
}

function getBoardStateString(roomState) {
  switch (roomState) {
    case "roomSelecting":
      return "Please, choose a card";
    case "roomShow":
      return "Result:";
  }
}

function getCardBoardState(state) {
  let inner = "<ul>";
  for (const player in state.Players) {
    switch (state.RoomState) {
      case "roomSelecting":
        inner += renderSelectionStep(state.Players[player]);
        break;
      case "roomShow":
        inner += renderSelectedCard(state.Players[player]);
    }
  }
  inner += `</ul>`;
  return inner;
}

function renderSelectionStep(player) {
  switch (player.PlayerState) {
    case "playerNoSelectedCard":
      return `<li class="player-list noselection">${player.Player}</li>`;
    case "playerSelected":
      return `<li class="player-list selection">${player.Player}</li>`;
  }
}

function renderSelectedCard(player) {
  return `<li class="player-list">${player.Player} : ${player.Card}</li>`;
}

function applyCardStyle(selectionDiv, state, playerID) {
  const className = "selected";
  allBtns = selectionDiv.getElementsByTagName("button");
  for (var i = 0; i < allBtns.length; i++) {
    allBtns[i].classList.remove(className);
  }

  selectedCard = document.getElementById(
    `card${state.Players[playerID].Card}`
  );
  if (selectedCard) {
    selectedCard.classList.add(className);
  }
}
