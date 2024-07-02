import './App.css';
import useFederatedComponent from "./helpers";
import React, {Fragment, lazy, Suspense} from "react";
import {useQuery} from "@apollo/client";
import {FetchShellConfigDocument, ShellNavigation} from "./__generated__/graphql";
import {Link, Route, Routes} from "react-router-dom";
import Home from "./components/Home";

const RemoteAppLoader = (props: any) => {
    const app: ShellNavigation = props.app;

    const {Component, isError} = useFederatedComponent({
        remoteUrl: 'http://localhost:3001/remoteEntry.js',
        moduleToLoad: './Counter',
        remoteName: 'example',
    });

    if (isError) return <div className="content">Error loading remote component</div>;

    return (
        <div>
            <Suspense>
                {Component ? <Component/> : <label>Not found</label>}
            </Suspense>
        </div>
    )
}

const App = () => {
    const {loading, error, data} = useQuery(FetchShellConfigDocument, {
        variables: {
            tenant: "abc"
        }
    })

    if (error) return <div className="content">Could not load configuration</div>

    if (loading) return <div className="content">Loading</div>

    return (
        <>
            <nav className="navbar">
                <ul className="navbar-nav">
                    <li className="nav-item">
                        <Link to="/">Home</Link>
                    </li>
                    {data?.shellConfiguration?.categories?.map((row) => (
                        <li key={row?.category} className="nav-item has-dropdown">
                            <a href="#">{row!.title}</a>
                            {row?.entries ?
                                <ul className="dropdown">
                                    {row?.entries?.map((row) => (
                                        <li key={row?.title} className="dropdown-item">
                                            <Link to={row!.module.path}>{row!.title}</Link>
                                        </li>
                                    ))}
                                </ul>
                                : null
                            }
                        </li>
                    ))}
                </ul>
            </nav>
            <div className="content">
                <Routes key="routes">
                    <Route key="home" path="/" element={<Home/>}/>
                    {data?.shellConfiguration?.categories?.map((row) => {
                        return row?.entries?.map(value => (
                            <Route key={row?.title} path={value?.module.path} element={<RemoteAppLoader app={value}/>}/>
                        ))
                    })}
                </Routes>
            </div>
        </>
    );
};

export default App;
