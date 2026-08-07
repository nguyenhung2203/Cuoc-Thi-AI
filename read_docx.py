import zipfile
import xml.etree.ElementTree as ET
import sys
import io

# Ensure utf-8 output for Vietnamese chars
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

def read_docx(filename):
    try:
        with zipfile.ZipFile(filename, 'r') as docx:
            xml_content = docx.read('word/document.xml')
            tree = ET.fromstring(xml_content)
            
            ns = {'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}
            
            paragraphs = []
            for p in tree.findall('.//w:p', ns):
                texts = [node.text for node in p.findall('.//w:t', ns) if node.text]
                if texts:
                    paragraphs.append(''.join(texts))
            return '\n'.join(paragraphs)
    except Exception as e:
        return str(e)

content = read_docx('D:\\thi\\Cuoc-Thi-AI\\De cuong thuyet minh chi tiet san pham du thi (1).docx')
with open('docx_content.txt', 'w', encoding='utf-8') as f:
    f.write(content)
