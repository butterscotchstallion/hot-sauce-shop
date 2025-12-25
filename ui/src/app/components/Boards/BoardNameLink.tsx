import {NavLink} from "react-router";

interface IBoardNameLinkProps {
    isOfficial: boolean;
    displayName: string;
    slug: string;
}

export function BoardNameLink({isOfficial, displayName, slug}: IBoardNameLinkProps) {
    return (
        <>
            <NavLink to={`/boards/${slug}#`}>
                {isOfficial && <i className="pi pi-verified mr-1"/>}
                {!isOfficial && <i className="pi pi-list mr-1"></i>}
                {displayName}
            </NavLink>
        </>
    )
}