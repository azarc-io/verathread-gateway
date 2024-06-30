import './App.css';
import useFederatedComponent from "./helpers";

const App = () => {
    const { Component, isError } = useFederatedComponent({
        remoteUrl: 'http://localhost:3001/remoteEntry.js',
        moduleToLoad: './Counter',
        remoteName: 'example',
    });

    if (isError) return <div>Error loading remote component</div>;

  return (
    <div className="content">
      <h1>Rsbuild with React</h1>
      <p>Start building amazing things with Rsbuild.</p>
        {Component ? <Component/> : null}
    </div>
  );
};

export default App;
