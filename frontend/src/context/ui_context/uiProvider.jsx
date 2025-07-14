import { useState } from "react";
import UiContext from "./uiContext";

const UiProvider = ({ children }) => {
    const [isFolderTreeOpen, setIsFolderTreeOpen] = useState(true);
    const [authenticated, setAuthenticated] = useState(false);

    return <UiContext.Provider value={{
        isFolderTreeOpen,
        setIsFolderTreeOpen,
        authenticated,
        setAuthenticated
    }}>
        {children}
    </UiContext.Provider>;
};

export default UiProvider;