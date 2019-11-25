import React from 'react';

import { Router } from "@reach/router"
import WatchView from "./views/watch"
import { store } from "./store"
import { StoreProvider } from 'easy-peasy'
import { RoomStore } from "./store/roomStore"

let Home = () => <div>Home</div>
let Dash = () => <div>Dash</div>

let Watch = ({ id }) => (
  <RoomStore.Provider initialData={{ id }}>
    <WatchView></WatchView>
  </RoomStore.Provider>
)

function App() {
  return (
    <StoreProvider store={store}>
      <Router>
        <Home path="/" />
        <Dash path="dashboard" />
        <Watch path="watch/:id" />
      </Router>
    </StoreProvider>
  );
}

export default App;
