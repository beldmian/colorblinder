import { createContext, useContext, useState, useMemo } from 'react';

export const AppContext = createContext({
    theme: 'light',
    setTheme: null,
});

type Props = {
    children: React.ReactNode
}
export function AppWrapper({ children }: Props) {
    let [theme, setTheme] = useState('light');
    let themeMemo = useMemo(() => ({theme, setTheme}), [theme])

    return (
        <AppContext.Provider value={themeMemo}>
            {useMemo(() => 
            (<>
                {children}
            </>), [])}
        </AppContext.Provider>
    );
}
