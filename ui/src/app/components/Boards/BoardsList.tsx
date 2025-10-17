import {IBoard} from "./IBoard.ts";
import {classNames} from "primereact/utils";
import {DataView} from 'primereact/dataview';
import {ReactElement} from "react";

interface IBoardsListProps {
    boards: IBoard[];
}

export const BoardsList = (props: IBoardsListProps) => {
    const boardTemplate = (board: IBoard, index: number) => {
        return (
            <div className="col-12" key={board.id}>
                <div
                    className={classNames('flex flex-column xl:flex-row xl:align-items-start p-4 gap-4', {'border-top-1 surface-border': index !== 0})}>
                    <img className="w-9 sm:w-16rem xl:w-10rem shadow-2 block xl:block mx-auto border-round"
                         src='/images/hot-pepper.png'
                         alt={board.displayName}/>
                    <div
                        className="flex flex-column sm:flex-row justify-content-between align-items-center xl:align-items-start flex-1 gap-4">
                        <div className="flex flex-column align-items-center sm:align-items-start gap-3">
                            <div className="text-2xl font-bold text-900">{board.displayName}</div>

                            <div className="flex align-items-center gap-3">
                                <span className="flex align-items-center gap-2">
                                    <i className="pi pi-tag"></i>
                                    <span className="font-semibold">Foods</span>
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    };
    const listTemplate = (items: IBoard[]) => {
        if (!items || items.length === 0) return null;
        const list: ReactElement[] = items.map((product, index) => {
            return boardTemplate(product, index);
        });
        return <div className="grid grid-nogutter">{list}</div>;
    };
    return (
        <>
            <DataView value={props.boards} listTemplate={listTemplate} emptyMessage={"No boards available."}/>
        </>
    )
}