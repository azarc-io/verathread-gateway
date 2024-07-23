import {ComponentType, lazy, LazyExoticComponent, useEffect, useState} from 'react';
import loadComponent from './componentLoaderFromWebpackContainer';
import {getOrLoadRemote} from './getOrLoadRemote';
import {IFederatedComponent} from './types';

type LazyComponent = LazyExoticComponent<ComponentType<any>>;

export default function useFederatedComponent({
                                                  remoteName,
                                                  remoteUrl,
                                                  moduleToLoad,
                                                  shareScope = 'default'
                                              }: IFederatedComponent) {
    const key = `${remoteUrl}-${remoteName}-${shareScope}-${moduleToLoad}`;

    const [Component, setComponent] = useState<LazyComponent | null>(null);
    const [isError, setIsError] = useState(false);

    useEffect(() => {
        console.log('effect 2')
        setComponent(null);
        setIsError(false);

        getOrLoadRemote({
            remoteName, remoteUrl, shareScope
        }).then(() => {
            const Comp = lazy(loadComponent({remoteName, moduleToLoad}));
            setComponent(Comp);
            setIsError(false);
            console.log('module loaded', Comp)
        }).catch((err) => {
            console.error(err);
            setIsError(true);
        })

        return () => {
            console.log('cleaning up module')
            setComponent(null)
            setIsError(false)
        }
    }, [key]);

    return {isError, Component};
}
