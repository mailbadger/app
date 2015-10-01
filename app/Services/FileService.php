<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 15.8.15
 * Time: 16:47
 */

namespace newsletters\Services;


use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Log;
use PHPExcel;
use PHPExcel_Exception;
use PHPExcel_IOFactory;
use PHPExcel_Worksheet;
use PHPExcel_Worksheet_Row;
use PHPExcel_Worksheet_RowIterator;

class FileService
{

    /**
     * Reads the excel file and returns an array of subscribers
     *
     * @param $file
     * @return Collection
     */
    public function importSubscribers($file)
    {
        $obj = $this->loadFile($file);

        return $this->readFile($obj);
    }

    /**
     * Export table data to excel file
     *
     * @param array $data
     * @param array $header
     * @param PHPExcel $excelObj
     * @return PHPExcel
     */
    public function exportData(array $data, array $header, PHPExcel $excelObj)
    {
        $i = 0;
        foreach ($header as $field) {
            $excelObj->getActiveSheet()->setCellValueByColumnAndRow($i++, 1, $field);
        }

        $j = 2;
        foreach ($data as $row) {
            $i = 0;
            foreach ($header as $field) {
                $value = (isset($row[$field])) ? $row[$field] : '';
                $excelObj->getActiveSheet()->setCellValueByColumnAndRow($i++, $j, $value);
            }
            $j++;
        }

        return $excelObj;
    }

    /**
     * Creates a new PHPExcel object
     *
     * @param $title
     * @param string $description
     * @return PHPExcel
     * @throws PHPExcel_Exception
     */
    public function createExcelFile($title, $description = '')
    {
        $excelObj = new PHPExcel();
        $excelObj->getProperties()->setTitle($title)->setDescription($description);
        $excelObj->setActiveSheetIndex(0);

        return $excelObj;
    }

    /**
     * Creates an excel writer
     *
     * @param PHPExcel $excelObj
     * @param $writerType
     * @return \PHPExcel_Writer_IWriter
     */
    public function createWriter(PHPExcel $excelObj, $writerType = 'CSV')
    {
        return PHPExcel_IOFactory::createWriter($excelObj, $writerType);
    }

    /**
     * @param $file
     * @return PHPExcel
     */
    public function loadFile($file)
    {
        $reader = PHPExcel_IOFactory::createReaderForFile($file);
        $reader->setReadDataOnly(true);

        return $reader->load($file);
    }

    /**
     * Returns array of all subscribers to be imported
     *
     * @param PHPExcel $object
     * @return Collection
     */
    public function readFile(PHPExcel $object)
    {
        $sheet = $this->getWorksheet($object);

        return $this->readRows($sheet);
    }

    /**
     * @param PHPExcel $object
     * @param int $sheetIndex
     * @return \PHPExcel_Worksheet
     */
    public function getWorksheet(PHPExcel $object, $sheetIndex = 0)
    {
        $object->setActiveSheetIndex($sheetIndex);

        return $object->getActiveSheet();
    }

    /**
     * Parse all rows from the file and return an array of subscribers
     *
     * @param PHPExcel_Worksheet $worksheet
     * @return Collection
     */
    public function readRows(PHPExcel_Worksheet $worksheet)
    {
        $subscribers = new Collection();

        try {
            $headerRow = $this->getHeaderRow($worksheet->getRowIterator(1, 1));

            foreach ($worksheet->getRowIterator(2) as $row) { //Iterate from the second row
                $subscriber = $this->getSubscriberData($row, $headerRow);
                $subscribers->push($subscriber);
            }

            return $subscribers;
        } catch (PHPExcel_Exception $e) {
            Log::error($e->getMessage());

            return $subscribers;
        }
    }

    /**
     * Get the first row values
     *
     * @param PHPExcel_Worksheet_RowIterator $rowIterator
     * @return array
     */
    public function getHeaderRow(PHPExcel_Worksheet_RowIterator $rowIterator)
    {
        $cells = $rowIterator->current()->getCellIterator();
        $cells->setIterateOnlyExistingCells(false);

        $headerRow = [];
        foreach ($cells as $cell) {
            $headerRow[$cell->getColumn()] = $cell->getValue();
        }

        return $headerRow;
    }

    /**
     * Parse each cell in the row and return an array of subscriber and custom fields
     *
     * @param PHPExcel_Worksheet_Row $row
     * @param $headerRow
     * @return array
     */
    public function getSubscriberData(PHPExcel_Worksheet_Row $row, $headerRow)
    {
        $cells = $row->getCellIterator();
        $cells->setIterateOnlyExistingCells(false);

        $data = $subscriber = $customFields = [];

        foreach ($cells as $cell) {
            $column = $headerRow[$cell->getColumn()];

            if ('name' === strtolower($column) || 'email' === strtolower($column)) {
                $subscriber[$column] = $cell->getValue();
            } else {
                $customFields[] = ['name' => $column, 'value' => $cell->getValue()];
            }
        }

        $data['subscriber'] = $subscriber;
        $data['custom_fields'] = $customFields;

        return $data;
    }
}