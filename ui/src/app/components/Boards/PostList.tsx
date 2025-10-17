import {IBoardPost} from "./IBoardPost.ts";
import {ReactElement} from "react";
import {Card} from "primereact/card";
import {NavLink} from "react-router";
import TimeAgo from "react-timeago";

interface IPostListProps {
    posts: IBoardPost[];
}

export default function PostList({posts}: IPostListProps) {
    const postElement = (post: IBoardPost): ReactElement => {
        let postImagePath: string = '/images/hot-pepper.png';
        if (post.thumbnailFilename) {
            postImagePath = `/images/posts/${post.thumbnailFilename}`;
        }
        return (
            <Card key={`post-${post.id}`} title={post.title}>
                <div>
                    <div className="flex flex-column">
                        <img src={postImagePath} alt={post.title}/>
                        <div className="ml-4">
                            {post.postText}
                        </div>
                    </div>
                </div>
                <div>
                    <ul className="flex flex-wrap items-center text-gray-900 dark:text-white">
                        <li className="">
                            <NavLink to={`/boards/${post.boardSlug}`}>
                                <i className="pi pi-list mr-1"></i> {post.boardName}
                            </NavLink>
                        </li>
                        <li className="ml-4 mr-4">&bull;</li>
                        {/*<li>*/}
                        {/*    <NavLink to={`/users/${post.createdByUserSlug}`}>*/}
                        {/*        <i className="pi pi-user mr-1"></i> {post.createdByUsername}*/}
                        {/*    </NavLink>*/}
                        {/*</li>*/}
                        <li className="cursor-help">
                            <TimeAgo date={post.createdAt}/>
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