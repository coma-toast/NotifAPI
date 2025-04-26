// https://www.joshwcomeau.com/animation/3d-button/
export default function Button({}) {
    return (
        <>
            <button className="pushable">
                <span className="shadow"></span>
                <span className="edge"></span>
                <span className="front">{}</span>
            </button>
        </>
    );
}
