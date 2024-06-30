import './counter.module.scss'
import {useState} from "react";

function Counter() {
    const [curCount, setCurCount] = useState(0)

    return (
        <>
            <button onClick={() => setCurCount((count) => count + 1)}>
                remote count is {curCount}
            </button>
        </>
    )
}

export default Counter
