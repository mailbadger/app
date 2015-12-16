<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class Complaint extends Model implements Transformable
{
    use TransformableTrait;

    protected $fillable = [
        'recipient',
        'sender',
        'type',
        'timestamp', 
        'sent_email_id',
    ];

}
