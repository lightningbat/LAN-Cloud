import { useContext } from "react";
import ExplorerContext from "./explorerContext";

const useExplorerContext = () => useContext(ExplorerContext);

export default useExplorerContext;