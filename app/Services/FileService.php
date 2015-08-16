<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 15.8.15
 * Time: 16:47
 */

namespace newsletters\Services;


use Exception;
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
    public function importSubscribersFromFile($file)
    {
        try {
            $obj = $this->loadFile($file);

            return $this->readFile($obj);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return new Collection();
        }
    }

    /**
     * @param $file
     * @return PHPExcel|null
     */
    private function loadFile($file)
    {
        try {
            $reader = PHPExcel_IOFactory::createReaderForFile($file);
            $reader->setReadDataOnly(true);

            return $reader->load($file);
        } catch (PHPExcel_Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
    }

    /**
     * Returns array of all subscribers to be imported
     *
     * @param PHPExcel $object
     * @return Collection
     */
    private function readFile(PHPExcel $object)
    {
        $sheet = $this->getWorksheet($object);

        return $this->readRows($sheet);
    }

    /**
     * @param PHPExcel $object
     * @param int $sheetIndex
     * @return null|\PHPExcel_Worksheet
     */
    private function getWorksheet(PHPExcel $object, $sheetIndex = 0)
    {
        try {
            $object->setActiveSheetIndex($sheetIndex);

            return $object->getActiveSheet();
        } catch (PHPExcel_Exception $e) {
            return null;
        }
    }

    /**
     * Parse all rows from the file and return an array of subscribers
     *
     * @param PHPExcel_Worksheet $worksheet
     * @return Collection
     */
    private function readRows(PHPExcel_Worksheet $worksheet)
    {
        $subscribers = new Collection();

        try {
            $headerRow = $this->getHeaderRow($worksheet->getRowIterator(1, 1));

            foreach ($worksheet->getRowIterator(2) as $row) { //Iterate from the second row
                $subscriber = $this->readCells($row, $headerRow);
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
    private function getHeaderRow(PHPExcel_Worksheet_RowIterator $rowIterator)
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
    private function readCells(PHPExcel_Worksheet_Row $row, $headerRow)
    {
        try {
            $cells = $row->getCellIterator();
            $cells->setIterateOnlyExistingCells(false);

            $data = $subscriber = $customFields = [];

            foreach ($cells as $cell) {
                $column = $headerRow[$cell->getColumn()];

                if ($column == 'name' || $column == 'email') {
                    $subscriber[$column] = $cell->getValue();
                } else {
                    $customFields[$column] = $cell->getValue();
                }
            }

            $data['subscriber'] = $subscriber;
            $data['custom_fields'] = $customFields;

            return $data;
        } catch (PHPExcel_Exception $e) {
            Log::error($e->getMessage());

            return [];
        }
    }
}