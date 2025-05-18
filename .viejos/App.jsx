//DeBakatas (youtube)

import React from 'react';

const Component = () => {
    const [counter, setCounter] = React.useState(0);
    return (
        <button
            on click={() => setCounter(prev => prev + 1)}>
            I love React a {counter}!
        </button>
    );
};

export default Component;