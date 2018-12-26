function createStyle() {
    let style = document.createElement('link');
    style.type ='text/css';
    style.rel = 'stylesheet';
    style.href = 'https://cdn.jsdelivr.net/npm/aplayer@1.10.1/dist/APlayer.min.css';
    document.head.appendChild(style);
}

function createScript() {
    let script = document.createElement('script');
    script.setAttribute('src', 'https://cdn.jsdelivr.net/npm/aplayer@1.10.1/dist/APlayer.min.js');
    document.body.appendChild(script);
    document.body.removeChild(script);
}

function getLrc(songid) {
    this.url = `https://api.itswincer.com/cloudmusic/media?id=${songid}`;
    return fetch(url)
}
function getDetail(songId) {
    this.url = `https://api.itswincer.com/cloudmusic/detail?ids=[${songId}]`;
    return fetch(url);
}

function getMusicId() {
    let idElement = document.querySelector('#aplayer');
    return idElement.getAttribute('musicid');
}

function createPlayer(detail) {
    new APlayer({
        container: document.querySelector("#aplayer"),
        narrow: false,
        autoplay: false,
        showlrc: 3,
        mutex: true,
        lrcType: 1,
        theme: "#ad7a86",
        music: [{
            title: detail.songs[0].name,
            author: detail.songs[0].artists[0].name,
            url: `https://music.163.com/song/media/outer/url?id=${songId}.mp3`,
            pic: detail.songs[0].album.picUrl + '?param=130y130',
            lrc: detail.lrc,
        }]
    });
    document.querySelectorAll('.aplayer-lrc p').forEach(e => {e.style.fontSize="12px"});
}

const songId = getMusicId();
createStyle();createScript();

Promise.all([getLrc(songId), getDetail(songId)])
    .then(async (resps) => {
        let [lrc, detail] = [await resps[0].text(), await resps[1].text()];
        [lrc, detail] = [JSON.parse(lrc), JSON.parse(detail)]
        if (lrc.nolyric) {
            lrc['lyric'] = '[00:00.00]纯音乐，请欣赏\n';
        }
        detail['lrc'] = lrc.lyric
        createPlayer(detail)
    })
