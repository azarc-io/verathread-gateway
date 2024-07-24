import {ShellNavigation} from "../__generated__/graphql";
import React, {ComponentType, lazy, LazyExoticComponent, Suspense, useEffect, useState} from "react";
import {getOrLoadRemote} from "../helpers/getOrLoadRemote";
import loadComponent from "../helpers/componentLoaderFromWebpackContainer";

type LazyComponent = LazyExoticComponent<ComponentType<any>>;

const Loading = () => {
    return (
        <div>Loading</div>
    )
}

const Component = (props: any) => {
    const app: ShellNavigation = props.app;

    const [Component, setComponent] = useState<LazyComponent | null>(null);
    const [isError, setIsError] = useState(false);

    useEffect(() => {
        getOrLoadRemote({
            remoteName: app.module.moduleName,
            remoteUrl: app.module.remoteEntry,
            shareScope: 'default'
        }).then(() => {
            const Comp = lazy(loadComponent({
                remoteName: app.module.moduleName,
                moduleToLoad: app.module.exposedModule
            }));
            setComponent(Comp);
            setIsError(false);
        }).catch((err) => {
            console.error(err);
            setIsError(true);
        })

        return () => {
            setComponent(null)
            setIsError(false)
        }
    }, [app]);

    if (isError) return <div className="content">Error loading remote component</div>;

    return (
        <div>
            <Suspense fallback={<Loading/>}>
                {Component ? <Component /> : null}
            </Suspense>
        </div>
    )
}

const RemoteAppLoader = (props: any) => {
    const app: ShellNavigation = props.app;
    return (
        <Component key={app.id} app={app}/>
    )
}

export default RemoteAppLoader
