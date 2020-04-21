import React from "react";
import ReactDOM from "react-dom";
import { Provider } from "react-redux";
import App from "./App";
import store from "./store";

// import * as serviceWorker from './serviceWorker';

console.log(
  "---------------\n" +
    "Hey there, curious fellow! While we love inquisitive minds, please note that reverse-engineering any part of CanvasCBL is against our terms of service. Thank you for using CanvasCBL!\n\n" +
    "If you're a developer, you should totally check out our API at https://go.canvascbl.com/docs!\n" +
    "---------------"
);

ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
// serviceWorker.unregister();
