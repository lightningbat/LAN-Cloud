import './style.scss'
import { FolderIcon, FileIcon, CheckBoxCheckedIcon, CheckBoxUnCheckedIcon } from '../../../../../icons'
import { useExplorerContext } from '../../../../../context/explorer_context';

export default function Content() {
    const { loading, 
        foldersData, 
        filesData, 
        selectedFolderId, 
        setSelectedFolderId, 
        selectedTagState, 
        tagsItems, 
        selectionMode,
        selectedItems,
        updateSelectedItems } = useExplorerContext();


    const toggleSelect = (id) => {
        if (selectedItems[id]) {
            updateSelectedItems((draft) => {
                delete draft[id];
            });
        } else {
            updateSelectedItems((draft) => {
                draft[id] = "folder";
            });
        }
    }

    return (
        <div className="explorer-content">
            <div className="list">

                <div className="heading">
                    <p className={`name ${selectionMode ? "selection-mode" : ""}`}>Name</p>
                    <p className="size">Size</p>
                    <p className="date">Modified</p>
                </div>
            
                {selectedFolderId && loading !== selectedFolderId && Object.keys(foldersData[selectedFolderId]?.sub_folders).map((folder_id) =>
                    <Folder key={folder_id} onClick={() => setSelectedFolderId(folder_id)} selectionMode={selectionMode} selectedItems={selectedItems} toggleSelect={()=> toggleSelect(folder_id)} id={folder_id} name={foldersData[folder_id]?.name} size={foldersData[folder_id]?.size} modified={foldersData[folder_id]?.modified_time} />)
                }
                {selectedFolderId && loading !== selectedFolderId && Object.keys(foldersData[selectedFolderId]?.files).map((file_id) =>
                    <File key={file_id} selectionMode={selectionMode} selectedItems={selectedItems} toggleSelect={()=> toggleSelect(file_id)} id={file_id} name={filesData[file_id]?.name} size={filesData[file_id]?.size} modified={filesData[file_id]?.modified_time} />
                )}
                {selectedTagState && selectedTagState.type === "SystemTags" && Object.keys(tagsItems?.[selectedTagState.type]?.[selectedTagState?.id] || {}).map((file_id) =>
                    <File key={file_id} selectionMode={selectionMode} selectedItems={selectedItems} toggleSelect={()=> toggleSelect(file_id)} id={file_id} name={filesData[file_id]?.name} size={filesData[file_id]?.size} modified={filesData[file_id]?.modified_time} />
                )}
                {selectedTagState && selectedTagState.type === "UserTags" && <>
                    {
                        tagsItems.UserTags[selectedTagState.id]?.folders.map((folder_id) => {
                            const folder = foldersData[folder_id]
                            return (
                                <Folder key={folder_id} onClick={() => setSelectedFolderId(folder_id)} selectionMode={selectionMode} selectedItems={selectedItems} toggleSelect={()=> toggleSelect(folder_id)} id={folder_id} name={folder?.name} size={folder?.size} modified={folder?.modified_time} />
                            )
                        })
                    }
                    {
                        tagsItems.UserTags[selectedTagState.id]?.files.map((file_id) => {
                            const file = filesData[file_id]
                            return (
                                <File key={file_id} selectionMode={selectionMode} selectedItems={selectedItems} toggleSelect={()=> toggleSelect(file_id)} id={file_id} name={file?.name} size={file?.size} modified={file?.modified_time} />
                            )
                        })
                    }
                </>
                }
            </div>
        </div>
    )
}

function Folder({ id, name, size, modified, onClick, selectionMode, selectedItems, toggleSelect }) {
    if (name === undefined || size === undefined || modified === undefined) return null
    return (
        <div className='list-item'>
            {selectionMode && <div className='checkbox' onClick={toggleSelect}>
                {selectedItems[id] ? <CheckBoxCheckedIcon /> : <CheckBoxUnCheckedIcon />}
            </div>}
            <div className={`item-content ${selectedItems[id] ? "selected" : ""}`} onClick={onClick}>
                <Name icon={<FolderIcon style={{ width: '1rem', height: '1rem' }} />} name={name} />
                <p className="size">{sizeTranslator(size)}</p>
                <p className="date">{new Date(modified * 1000).toLocaleString()}</p>
            </div>
        </div>
    )
}

function File({ id, name, size, modified, selectionMode, selectedItems, toggleSelect }) {
    if (name === undefined || size === undefined || modified === undefined) return null
    return (
        <div className='list-item'>
            {selectionMode && <div className='checkbox' onClick={toggleSelect}>
                {selectedItems[id] ? <CheckBoxCheckedIcon /> : <CheckBoxUnCheckedIcon />}
            </div>}
            <div className={`item-content ${selectedItems[id] ? "selected" : ""}`}>
                <Name icon={<FileIcon style={{ width: '1rem', height: '1rem' }} />} name={name} />
                <p className="size">{sizeTranslator(size)}</p>
                <p className="date">{new Date(modified * 1000).toLocaleString()}</p>
            </div>
        </div>
    )
}

function Name({ icon, name }) {
    return (
        <div className='name'>
            {icon}
            <p>{name}</p>
        </div>
    )
}

function sizeTranslator(size) {
    if (size < 1024) return `${size} B`
    if (size >= 1024 && size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
    if (size >= 1024 * 1024 && size < 1024 * 1024 * 1024) return `${(size / 1024 / 1024).toFixed(2)} MB`
    if (size >= 1024 * 1024 * 1024) return `${(size / 1024 / 1024 / 1024).toFixed(2)} GB`
}