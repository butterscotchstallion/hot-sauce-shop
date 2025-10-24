import dayjs from "dayjs";
import {Card} from "primereact/card";
import {NavLink} from "react-router";
import TimeAgo from "react-timeago";
import {IBoardPost} from "./IBoardPost.ts";
import {Button} from "primereact/button";
import {useEffect, useRef, useState} from "react";
import {Subject} from "rxjs";
import {addUpdateVote} from "./VoteService.ts";
import {VoteValue} from "./VoteValue.ts";
import {getPostDetail} from "./BoardsService.ts";

interface IBoardPostProps {
    boardPost: IBoardPost;
    voteMap: Map<number, number>;
    replyMap: Map<number, number>;
}

export default function BoardPost({boardPost, voteMap, replyMap}: IBoardPostProps) {
    const [createdAtFormatted, setCreatedAtFormatted] = useState<string>(dayjs(boardPost.createdAt).format('MMMM D, YYYY'));
    const [hasUpVoted, setHasUpVoted] = useState<boolean>(false);
    const [hasDownVoted, setHasDownVoted] = useState<boolean>(false);
    const vote$ = useRef<Subject<number>>(null);
    const post$ = useRef<Subject<IBoardPost>>(null);
    const postImagePath = useRef<string>('/images/hot-pepper.png');
    const [post, setPost] = useState<IBoardPost>(boardPost);

    const onUpvoteClicked = () => {
        setHasUpVoted(true);
        setHasDownVoted(false);
        vote(VoteValue.Upvote);
    }

    const onDownvoteClicked = () => {
        setHasDownVoted(true);
        setHasUpVoted(false);
        vote(VoteValue.Downvote);
    }

    const vote = (voteValue: number) => {
        vote$.current = addUpdateVote(post.id, voteValue);
        vote$.current.subscribe({
            next: (_) => {
                post$.current = getPostDetail(post.boardSlug, post.slug);
                post$.current.subscribe({
                    next: (updatedPost: IBoardPost) => {
                        console.info(`Updated post #${post.id} with vote value ${voteValue}: ${JSON.stringify(updatedPost, null, 2)}`);
                        setPost(updatedPost);
                    },
                    error: (err: string) => {
                        console.log('Error getting post: ' + err);
                    }
                })
            },
            error: (err: string) => {
                setHasUpVoted(false);
                setHasDownVoted(false);
                console.log('Error voting: ' + err);
            }
        });
    }

    const header = () => {
        return (
            <i className="pi pi-ellipsis-h cursor-pointer float-right inline-block mt-4 mr-4"></i>
        )
    }

    useEffect(() => {
        setCreatedAtFormatted(dayjs(post.createdAt).format('MMMM D, YYYY'));
    }, [post]);

    /**
     * Updates posts on load
     */
    useEffect(() => {
        setPost(boardPost);

        if (boardPost.thumbnailFilename) {
            postImagePath.current = `/images/posts/${boardPost.thumbnailFilename}`;
        }

        if (voteMap.has(post.id)) {
            const voteValue = voteMap.get(post.id);
            if (voteValue === VoteValue.Upvote) {
                setHasUpVoted(true);
            } else if (voteValue === VoteValue.Downvote) {
                setHasDownVoted(true);
            }
        }

        return () => {
            if (vote$.current) {
                vote$.current.unsubscribe();
            }
            if (post$.current) {
                post$.current.unsubscribe();
            }
        }
    }, [boardPost, voteMap]);

    return (
        <Card key={`post-${post.id}`}
              header={header}
              title={<NavLink to={`/boards/${post.boardSlug}/posts/${post.slug}`}>{post.title}</NavLink>}
              className="mb-4">
            <div>
                <div className="flex flex-column">
                    <img src={postImagePath.current} alt={post.title}/>
                    <div className="ml-6 min-h-[2rem]">
                        {post.postText}
                    </div>
                </div>
            </div>
            <div className="mt-6">
                <ul className="flex flex-wrap items-center text-gray-900 dark:text-white">
                    <li className="inline-block w-[120px]">
                        <div className="w-full flex justify-between items-center">
                            <Button
                                icon={hasUpVoted ? 'pi pi-thumbs-up-fill' : 'pi pi-thumbs-up'}
                                title="This is a high-quality post"
                                disabled={hasUpVoted}
                                onClick={() => onUpvoteClicked()}
                            />
                            <strong>{post.voteSum}</strong>
                            <Button
                                icon={hasDownVoted ? 'pi pi-thumbs-down-fill' : 'pi pi-thumbs-down'}
                                title="This is a low-quality post"
                                disabled={hasDownVoted}
                                onClick={() => onDownvoteClicked()}
                            />
                        </div>
                    </li>
                    <li className="ml-4 mr-4">&bull;</li>
                    <li>
                        <NavLink to={`/boards/${post.boardSlug}/posts/${post.slug}`}>
                            <i className="pi pi-reply"></i> {replyMap?.get(post.id) || 0}
                        </NavLink>
                    </li>
                    <li className="ml-4 mr-4">&bull;</li>
                    <li>
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
}