import './style.scss'
import { FolderIcon, FileIcon } from '../../../../../icons'
import { useExplorerContext } from '../../../../../context/explorer_context';

export default function Content() {
    const { loading, foldersData, filesData, selectedFolderId, setSelectedFolderId, selectedTagState, tagsItems } = useExplorerContext();

    return (
        <div className="explorer-content">
            <div className="heading">
                <p className="name">Name</p>
                <p className="size">Size</p>
                <p className="date">Modified</p>
            </div>

            <div className="list">
                {selectedFolderId && loading !== selectedFolderId && Object.keys(foldersData[selectedFolderId]?.sub_folders).map((folder_id) => 
                    <Folder key={folder_id} onClick={() => setSelectedFolderId(folder_id)} name={foldersData[folder_id]?.name} size={foldersData[folder_id]?.size} modified={foldersData[folder_id]?.modified_time} />)
                }
                {selectedFolderId && loading !== selectedFolderId && Object.keys(foldersData[selectedFolderId]?.files).map((file_id) =>
                    <File key={file_id} name={filesData[file_id]?.name} size={filesData[file_id]?.size} modified={filesData[file_id]?.modified_time} />
                )}
                {selectedTagState && Object.keys(tagsItems?.[selectedTagState.type]?.[selectedTagState?.id] || {}).map((file_id) =>
                    <File key={file_id} name={filesData[file_id]?.name} size={filesData[file_id]?.size} modified={filesData[file_id]?.modified_time} />
                )}
            </div>
        </div>
    )
}

function Folder({ name, size, modified, onClick }) {
    if (name === undefined || size === undefined || modified === undefined) return null
    return (
        <div className='list-item' onClick={onClick}>
            <Name icon={<FolderIcon style={{ width: '1rem', height: '1rem' }} />} name={name} />
            <p className="size">{sizeTranslator(size)}</p>
            <p className="date">{new Date(modified*1000).toLocaleString()}</p>
        </div>
    )
}

function File({ name, size, modified }) {
    if (name === undefined || size === undefined || modified === undefined) return null
    return (
        <div className='list-item'>
            <Name icon={<FileIcon style={{ width: '1rem', height: '1rem' }} />} name={name} />
            <p className="size">{sizeTranslator(size)}</p>
            <p className="date">{new Date(modified*1000).toLocaleString()}</p>
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