import {ShellNavigation} from "../__generated__/graphql";
import useFederatedComponent from "../helpers";
import React, {Suspense} from "react";
import ComponentNotFound from "./ComponentNotFound";

const RemoteAppLoader = (props: any) => {
    const app: ShellNavigation = props.app;
    console.log('loading', app)
    const {Component, isError} = useFederatedComponent({
        remoteUrl: 'http://localhost:3001/remoteEntry.js',
        moduleToLoad: './Counter',
        remoteName: 'example',
    });

    if (isError) return <div className="content">Error loading remote component</div>;

    return (
        <div>
            <Suspense>
                {Component ? <Component/> : <ComponentNotFound/>}
            </Suspense>
        </div>
    )
}

export default RemoteAppLoader
