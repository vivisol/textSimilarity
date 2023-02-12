# 这是一个demo，用来测试jina和语义分析工具similarity的

from jina import DocumentArray, Document, Executor, Flow, requests
from similarities import Similarity


class check(Executor):
    @requests
    async def add_text(self, docs: DocumentArray, **kwargs):
        print("checkSimilarity service is been called")
        m = Similarity()
        print(docs[0].text, docs[1].text)

        r = m.similarity(docs[0].text, docs[1].text)
        t = str(r.numpy()[0][0])
      
        resp = DocumentArray([Document(text=t)])
        print("score:",t)
       
        return resp

f = Flow(port=12345).add(uses=check, replicas=1)

if __name__ == '__main__':
    with f:
        f.block()