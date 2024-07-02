import {ShellNavigation} from "../__generated__/graphql";
import useFederatedComponent from "../helpers";
import React, {lazy, Suspense} from "react";
import {useParams} from "react-router-dom";

// export default function RemoteAppLoader() {
//     const {app} = useParams();
//     console.log('loading', app)
//
//     const {Component, isError} = useFederatedComponent({
//         remoteUrl: 'http://localhost:3001/remoteEntry.js',
//         moduleToLoad: './Counter',
//         remoteName: 'example',
//     });
//
//     if (isError) return <div className="content">Error loading remote component</div>;
//
//     return (
//         <div>
//             <Suspense>
//                 {Component ? <Component/> : <label>Not found</label>}
//             </Suspense>
//         </div>
//     )
// }

const RemoteLoader = lazy(() => {
    const {app} = useParams();
    console.log('loading', app)

    const {Component, isError} = useFederatedComponent({
        remoteUrl: 'http://localhost:3001/remoteEntry.js',
        moduleToLoad: './Counter',
        remoteName: 'example',
    });

    if (isError || !Component) {
        return import('./ComponentNotFound')
    }

    // return Component
    return import('./ComponentNotFound')
})

export default RemoteLoader
