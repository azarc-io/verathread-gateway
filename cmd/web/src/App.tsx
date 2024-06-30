import React, {Suspense} from 'react';
import logo from './logo.svg';
import './App.css';
import useFederatedComponent from 'mf-cra';

const remoteApps = [
  {
    "remoteUrl": "http://localhost:3001/remoteEntry.js",
    "moduleToLoad": "./ExampleApp",
    "remoteName": "example",
    "localRoute": "example1"
  }
]

function RemoteApp({ app }) {
  const { Component: RemoteComponent } = useFederatedComponent(app);

  return (
      <Suspense fallback='loading...'>
        {RemoteComponent && <RemoteComponent />}
      </Suspense>
  );
}

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <RemoteApp app={remoteApps[0]}/>
      </header>
    </div>
  );
}

export default App;
