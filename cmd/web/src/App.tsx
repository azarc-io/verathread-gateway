import './App.css';
import {useSubscription} from "@apollo/client";
import {ShellConfigEventType, SubscribeToShellConfigDocument} from "./__generated__/graphql";
import {Link, Route, Routes} from "react-router-dom";
import Home from "./components/Home";
import RemoteAppLoader from "./components/RemoteLoader";

const App = () => {
    const {loading, error, data} = useSubscription(SubscribeToShellConfigDocument, {
        variables: {
            tenant: "abc",
            events: [
                ShellConfigEventType.Initial,
                ShellConfigEventType.Rebuild,
                ShellConfigEventType.Updated
            ]
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
                    {data?.shellConfiguration?.configuration.categories?.map((row) => (
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
                    {data?.shellConfiguration?.configuration.categories?.map((row) => {
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
