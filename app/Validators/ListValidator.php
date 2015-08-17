<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 17.8.15
 * Time: 21:19
 */

namespace newsletters\Validators;


use Illuminate\Validation\Validator;
use newsletters\Services\FileService;

class ListValidator
{

    /**
     * @var FileService
     */
    private $fileService;

    public function __construct(FileService $fileService)
    {
        $this->fileService = $fileService;
    }

    public function validateCheckFields($attribute, $value, $parameters, Validator $validator)
    {
        $file = array_get($validator->getFiles(), $parameters[0]);
        $listId = array_get($validator->getData(), $parameters[1]);

        $obj = $this->fileService->loadFile($file);
        $sheet = $this->fileService->getWorksheet($obj, 0);

        $headerRow = $this->fileService->getHeaderRow($sheet->getRowIterator(1, 1));

        unset($file);

        return false;
    }
}