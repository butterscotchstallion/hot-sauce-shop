import {IBoardPost} from "./IBoardPost.ts";
import {ReactElement} from "react";
import {Card} from "primereact/card";
import {NavLink} from "react-router";
import TimeAgo from "react-timeago";
import dayjs from "dayjs";

interface IPostListProps {
    posts: IBoardPost[];
}

export default function PostList({posts}: IPostListProps) {
    const postElement = (post: IBoardPost): ReactElement => {
        const createdAtFormatted = dayjs(post.createdAt).format('MMM DD, YYYY hh:mm A');
        let postImagePath: string = '/images/hot-pepper.png';
        if (post.thumbnailFilename) {
            postImagePath = `/images/posts/${post.thumbnailFilename}`;
        }
        return (
            <Card key={`post-${post.id}`} title={post.title} className="mb-4">
                <div>
                    <div className="flex flex-column">
                        <img src={postImagePath} alt={post.title}/>
                        <div className="ml-6 min-h-[2rem]">
                            {post.postText}
                        </div>
                    </div>
                </div>
                <div className="mt-6">
                    <ul className="flex flex-wrap items-center text-gray-900 dark:text-white">
                        <li className="">
                            <NavLink to={`/boards/${post.boardSlug}`}>
                                <i className="pi pi-list mr-1"></i> {post.boardName}
                            </NavLink>
                        </li>
                        <li className="ml-4 mr-4">&bull;</li>
                        <li>
                            <NavLink to={`/users/${post.createdByUserSlug}`}>
                                <i className="pi pi-user mr-1"></i> {post.createdByUsername}
                            </NavLink>
                        </li>
                        <li className="ml-4 mr-4">&bull;</li>
                        <li className="cursor-help">
                            <TimeAgo date={post.createdAt} title={createdAtFormatted}/>
                        </li>
                    </ul>
                </div>
            </Card>
        )
    };
    const postList = posts?.map((post: IBoardPost): ReactElement => {
        return postElement(post);
    });
    return (
        <>
            {posts?.length > 0 && (
                <section className="mt-4">
                    {postList}
                </section>
            )}
        </>
    )
}