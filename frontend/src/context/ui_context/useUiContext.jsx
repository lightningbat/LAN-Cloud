import { useContext } from "react";
import UiContext from "./uiContext";

const useUiContext = () => useContext(UiContext);

export default useUiContext;