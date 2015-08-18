<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 17.8.15
 * Time: 21:19
 */

namespace newsletters\Validators;


use Exception;
use Illuminate\Validation\Validator;
use newsletters\Services\FieldService;
use newsletters\Services\FileService;

class ListValidator
{

    /**
     * @var FileService
     */
    private $fileService;

    /**
     * @var FieldService
     */
    private $fieldService;

    public function __construct(FileService $fileService, FieldService $fieldService)
    {
        $this->fileService = $fileService;
        $this->fieldService = $fieldService;
    }

    public function validateCheckFields($attribute, $value, $parameters, Validator $validator)
    {
        $file = array_get($validator->getFiles(), $parameters[0]);
        $listId = array_get($validator->getData(), $parameters[1]);

        try {
            $obj = $this->fileService->loadFile($file);
            $sheet = $this->fileService->getWorksheet($obj, 0);

            $headerRow = $this->fileService->getHeaderRow($sheet->getRowIterator(1, 1));

            if (empty($headerRow)) {
                return false;
            }

            $fields = $this->fieldService->findFieldsByListId($listId);
            $fields->push(['name' => 'name'])->push(['name' => 'email']);

            if (count($headerRow) !== $fields->count()) {
                return false;
            }

            $flag = false;
            foreach ($fields as $field) {
                $name = strtolower($field['name']);
                if (!($flag = in_array($name, $headerRow, true))) {
                    break;
                }
            }

            unset($file);

            return $flag;
        } catch (Exception $e) {
            return false;
        }
    }
}