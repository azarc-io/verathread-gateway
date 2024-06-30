import {lazy, useState} from 'react'
import { Route, Routes } from "react-router-dom";
import {useQuery} from "@apollo/client";
import {FetchShellConfigDocument} from "./__generated__/graphql.ts";

// import {
//   __federation_method_getRemote,
//   __federation_method_setRemote,
//   // @ts-ignore
// } from "__federation__";

import './App.scss'
import Header from "./components/Header.tsx";
import Loading from "./components/Loading.tsx";

// const DynamicRemoteApp = lazy(() => {
//   // values like { 'http://localhost:9000/assets/remoteEntry.js', 'remoteA', './RemoteARoot' }
//   const {url, name, module } = {
//     module: './RemoteARoot',
//     name: 'remoteA',
//     url: 'http://localhost:6110/module/remoteEntry.js'
//   };
//
//   __federation_method_setRemote(name, {
//     url: () => Promise.resolve(url),
//     format: "esm",
//     from: "vite",
//   });
//
//   const fm = __federation_method_getRemote(name, module)
//
//   console.log('default', fm)
//
//   return fm;
// });

function App() {
  const { loading, error, data } = useQuery(FetchShellConfigDocument, {
    variables: {
      tenant: "abc"
    }
  })

  const isDev = import.meta.env.DEV;
  const mode = import.meta.env.MODE;

  console.log('render', loading, error, data)

  return (
      <>
        <Header></Header>
        {loading ? (
            <Loading></Loading>
        ) : (
            <div>
              Loaded
              {/*<Routes>*/}
              {/*  <Route path="/" element={<DynamicRemoteApp />} />*/}
              {/*</Routes>*/}
            </div>
        )}

        <p>DEV variable: {JSON.stringify(isDev)}</p>
        <p>MODE variable: {mode}</p>
      </>
  )
}

export default App
